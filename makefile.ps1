$COMMAND = $args[0]

$NAME = "camera-services"
$OWNER = "byuoitav"
$PKG = "github.com/$OWNER/$NAME"
$DOCKER_URL = "docker.pkg.github.com"
$DOCKER_PKG = "$DOCKER_URL/$OWNER/$NAME"

Write-Output "PKG: $PKG"
Write-Output "DOCKER_PKG: $DOCKER_PKG"

$PRD_TAG_REGEX = "v[0-9]+\.[0-9]+\.[0-9]+"
$DEV_TAG_REGEX = "v[0-9]+\.[0-9]+\.[0-9]+-.+"


$COMMIT_HASH = Invoke-Expression "git rev-parse --short HEAD"
$TAG = Invoke-Expression "git rev-parse --short HEAD"
try {
    $NEW_TAG = Invoke-Expression "git describe --exact-match --tags HEAD"
    Write-Output "NEW_TAG: $NEW_TAG.Length"
    if ($NEW_TAG.Length -gt 0) {
        $TAG = $NEW_TAG
        Write-Output "The repo contains a tag: $TAG"
    }
}
catch {
    Write-Output "The repo does not contain a tag"
}

Write-Output "The TAG is: $TAG"

# go stuff
$PKG_LIST = Invoke-Expression "go list $PKG/..."
Write-Output "PKG_LIST: $PKG_LIST"


function All {
    Write-Output "All"
}

function Test {
    Write-Output "Test"
    Invoke-Expression "go test -v $PKG_LIST"
}

function Test-cov {
    Write-Output "Test-cov"
    Invoke-Expression "go test -coverprofile=coverage.txt -covermode=atomic $PKG_LIST"
}

function Lint {
    Write-Output "Lint"
    Invoke-Expression "golangci-lint run --test=false"
}

function Deps {
    Write-Output "Downloading Backend Dependencies"
    Invoke-Expression "go mod download"

    Write-Output "Downloading control frontend dependencies"
    Set-Location "cmd/control/web/"
    Invoke-Expression "npm install"
    Invoke-Expression "cd .."
    Invoke-Expression "cd .."
    Invoke-Expression "cd .."
    Write-Output "Exiting cmd/control/web/"

    Write-Output "Downloading spyglass frontend dependencies"
    Set-Location "cmd/spyglass/web/"
    Invoke-Expression "npm install"
    Invoke-Expression "cd .."
    Invoke-Expression "cd .."
    Invoke-Expression "cd .."
    Write-Output "Exiting cmd/spyglass/web/"
}

function Build {
    Write-Output "*******************Build Start**********************"

    # Create directories for compiled code
    New-Item -Path dist -ItemType Directory
    New-Item -Path dist/control -ItemType Directory
    New-Item -Path dist/spyglass -ItemType Directory

    # - Set default env vars for Windows if they dont exist
    if ($null -eq $env:CGO_ENABLED) { $env:CGO_ENABLED = 0 }
    if ($null -eq $env:GOOS) { $env:GOOS = "windows" }
    if ($null -eq $env:GOARCH) { $env:GOARCH = "amd64" }

    #Get current env vars
    $Start_CGO_ENABLED = $env:CGO_ENABLED
    Write-Output "Got CGO_ENABLED: $Start_CGO_ENABLED"

    $Start_GOOS = $env:GOOS
    Write-Output "Got GOOS: $Start_GOOS"

    $Start_GOARCH = $env:GOARCH
    Write-Output "Got GOARCH: $Start_GOARCH"

    # Set temp environment vars - same for all build actions
    Write-Output "Setting CGO_ENABLED, GOOS, and GOARCH for all build actions"
    Set-Item -Path env:CGO_ENABLED -Value 0
    Set-Item -Path env:GOOS -Value "linux"
    Set-Item -Path env:GOARCH -Value "amd64"

    # Build AVER for Linux AMD64
    Write-Output "Building AVER for Linux AMD64"
    Invoke-Expression "go build -v -o dist/aver-linux-amd64 ./cmd/aver/..."

    # Build AXIS for Linux AMD64
    Write-Output "Building AXIS for Linux AMD64"
    Invoke-Expression "go build -v -o dist/axis-linux-amd64 ./cmd/axis/..."

    # Build Control Backend for Linux AMD64
    Write-Output "Building Control Backend for Linux AMD64"
    Invoke-Expression "go build -v -o dist/control-linux-amd64 ./cmd/control/..."

    # Build Spyglass Backend for Linux AMD64
    Write-Output "Building Spyglass Backend for Linux AMD64"
    Invoke-Expression "go build -v -o dist/spyglass-linux-amd64 ./cmd/spyglass/..."


    # Build Control Frontend
    Write-Output "Building Control Frontend"
    Invoke-Expression "npm --prefix ./cmd/control/web run-script build"
    Write-Output "Moving files to  ./dist/control"
    Move-Item "./cmd/control/web/dist/" -Destination "./dist/control"

    # # Build Spyglass Frontend
    Write-Output "Building Spyglass Frontend"
    Invoke-Expression "npm --prefix ./cmd/spyglass/web run-script build"
    Write-Output "Moving files to  ./dist/spyglass"
    Move-Item "./cmd/spyglass/web/dist/" -Destination "./dist/spyglass"

    #Cleanup Env Vars
    Write-Output "*******************Build End**********************"
    Write-Output "Resetting env vars to start values"
    Set-Item -Path env:CGO_ENABLED -Value $Start_CGO_ENABLED
    Set-Item -Path env:GOOS -Value $Start_GOOS
    Set-Item -Path env:GOARCH -Value $Start_GOARCH
}

function Cleanup {
    Write-Output "Clean"
    Invoke-Expression "go clean"
    if (Test-Path -Path "dist") {
        Remove-Item dist -recurse
        Write-Output "Recursively deleted dist/"
    }
    else {
        Write-Output "No dist directory to delete"
    }
    if (Test-Path -Path "analog/dist") {
        Remove-Item analog/dist -recurse
        Write-Output "Recursively deleted dist/"
    }
    else {
        Write-Output "No analog/dist directory to delete"
    }
}

function DockerFunc {
    #can not just be docker because it creates an infinite loop
    $DevTag = ""
    $FileTag = ""

    Write-Output "Function Docker      Commit Hash: $COMMIT_HASH     Tag: $TAG"
    if ($COMMIT_HASH -eq $TAG) {
        
        $FileTag = $COMMIT_HASH
        $DevTag = "-dev"
        Write-Output "Building dev containers with tag $FileTag"
    }
    elseif ($TAG -match $DEV_TAG_REGEX) {
        $FileTag = $TAG
        $DevTag = "-dev"
        Write-Output "Building dev containers with tag $FileTag"
    }
    elseif ($TAG -match $PRD_TAG_REGEX) {
        $FileTag = $TAG
        $DevTag = ""
        Write-Output "Building prd containers with tag $FileTag"

    }
    else {
        Write-Output "Docker function quit unexpectedly. Commit Hash: $COMMIT_HASH     Tag: $TAG"
        return
    }

    Write-Output "Building container $DOCKER_PKG/aver${DevTag}:${FileTag}"
    Invoke-Expression "docker build -f dockerfile --build-arg NAME=aver-linux-amd64 -t $DOCKER_PKG/aver${DevTag}:${FileTag} dist"

    Write-Output "Building container $DOCKER_PKG/axis${DevTag}:${FileTag}"
    Invoke-Expression "docker build -f dockerfile --build-arg NAME=axis-linux-amd64 -t $DOCKER_PKG/axis${DevTag}:${FileTag} dist"

    Write-Output "Building container $DOCKER_PKG/control${DevTag}:${FileTag}"
    Invoke-Expression "docker build -f dockerfile-control --build-arg NAME=control-linux-amd64 -t $DOCKER_PKG/control${DevTag}:${FileTag} dist"

    Write-Output "Building container $DOCKER_PKG/camera-spyglass${DevTag}:${FileTag}"
    Invoke-Expression "docker build -f dockerfile-spyglass --build-arg NAME=spyglass-linux-amd64 -t $DOCKER_PKG/camera-spyglass${DevTag}:${FileTag} dist"
}

function Deploy {
    $DevTag = ""
    $FileTag = ""
    Write-Output "Deploy      Commit Hash: $COMMIT_HASH     Tag: $TAG"

    Write-Output "Logging into repo"    
    Invoke-Expression "docker login $DOCKER_URL -u $Env:DOCKER_USERNAME -p $Env:DOCKER_PASSWORD"
    
    if ($COMMIT_HASH -eq $TAG) {
        $FileTag = $COMMIT_HASH
        $DevTag = "-dev"
        Write-Output "Pushing dev containers with tag $FileTag"
    }
    elseif ($TAG -match $DEV_TAG_REGEX) {
        $FileTag = $TAG
        $DevTag = "-dev"
        Write-Output "Pushing dev containers with tag $FileTag"
    }
    elseif ($TAG -match $PRD_TAG_REGEX) {
        $FileTag = $TAG
        $DevTag = ""
        Write-Output "Pushing prd containers with tag $FileTag"
    }
    else {
        Write-Output "Deploy function quit unexpectedly. Commit Hash: $COMMIT_HASH     Tag: $TAG"
    }

    Write-Output "Pushing container $DOCKER_PKG/aver${DevTag}:${FileTag}"
    Invoke-Expression "docker push $DOCKER_PKG/aver${DevTag}:${FileTag}"

    Write-Output "Pushing container $DOCKER_PKG/axis${DevTag}:${FileTag}"
    Invoke-Expression "docker push $DOCKER_PKG/axis${DevTag}:${FileTag}"

    Write-Output "Pushing container $DOCKER_PKG/control${DevTag}:${FileTag}"
    Invoke-Expression "docker push $DOCKER_PKG/control${DevTag}:${FileTag}"

    Write-Output "Pushing container $DOCKER_PKG/camera-spyglass${DevTag}:${FileTag}"
    Invoke-Expression "docker push $DOCKER_PKG/camera-spyglass${DevTag}:${FileTag}"
}


if ($COMMAND -eq "All") {
    Cleanup
    Build
    All     
}
elseif ($COMMAND -eq "Test") {
    Deps
    Test
}
elseif ($COMMAND -eq "Test-cov") {
    Deps
    Test-cov
}
elseif ($COMMAND -eq "Lint") {
    Deps
    Lint
}
elseif ($COMMAND -eq "Deps") {
    Deps
}
elseif ($COMMAND -eq "Build") {
    Cleanup
    Deps
    Build
}
elseif ($COMMAND -eq "Clean") {
    Cleanup
}
elseif ($COMMAND -eq "Docker" ) {
    Cleanup
    Deps
    Build
    DockerFunc
}
elseif ($COMMAND -eq "Deploy" ) {
    Cleanup
    Deps
    Build
    DockerFunc
    Deploy
}
else {
    Write-Output "Please include a valid command parameter"
}
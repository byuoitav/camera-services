import {Component, HostListener, ViewChild, ElementRef, OnInit, EventEmitter, OnDestroy} from '@angular/core';
import {Router, ActivatedRoute} from "@angular/router";
import {Config, Camera, CameraPreset} from '../../../objects/objects';
import {HttpClient} from '@angular/common/http';
import { stringify } from '@angular/compiler/src/util';

@Component({
  selector: 'app-camera-feed',
  templateUrl: './camera-feed.component.html',
  styleUrls: ['./camera-feed.component.scss']
})
export class CameraFeedComponent implements OnInit, OnDestroy {
  rowHeight = "4:1.75"
  timeout = 0
  cameras: Config

  tilting = false
  zooming = false

  @ViewChild('stream') img:ElementRef;

  constructor(
    private router: Router,
    private http: HttpClient,
    public route: ActivatedRoute,
  ) {
    this.route.data.subscribe(data => {
      this.cameras = data.uiConfig;
      let title = this.cameras[0].displayName;
      let splitted = title.split(" ");
      title = splitted[0]+ " " + splitted[1];
      document.title = title + " Camera Control";
    })
  }
  ngOnInit() {
    setInterval(() => {
      this.timeout++
      if (this.timeout == 5) {
        console.log("preview timing out")
      }
    }, 1000)
  }

  ngOnDestroy() {
    console.log("ending http request")
    this.img.nativeElement.src = "";
  }


  @HostListener("window:resize", ["$event"])
  onResize(event) {
    if (window.innerWidth >= 1024 && window.innerHeight >= 768 && window.innerHeight <= 1056) {
      this.rowHeight = "4:2.5"
    } else {
      this.rowHeight = "4:1.75"
    }
  }

  exitRoom() {
    console.log("exiting room")
    this.router.navigate([""])
  }

  tiltUp = (cam: Camera) => {
    this.tilting = true
    this.timeout = 0
    console.log("tilting up", cam.tiltUp);
    if (!cam.tiltUp) {
      return;
    }

    this.http.get(cam.tiltUp).subscribe(resp => {
      console.log("resp", resp);
    }, err => {
      console.warn("err", err);
    });
  }

  tiltDown = (cam: Camera) => {
    this.tilting = true
    this.timeout = 0
    console.log("tilting down", cam.tiltDown);
    if (!cam.tiltDown) {
      return;
    }

    this.http.get(cam.tiltDown).subscribe(resp => {
      console.log("resp", resp);
    }, err => {
      console.warn("err", err);
    });
  }

  panLeft = (cam: Camera) => {
    this.tilting = true
    this.timeout = 0
    console.log("panning left", cam.panLeft);
    if (!cam.panLeft) {
      return;
    }

    this.http.get(cam.panLeft).subscribe(resp => {
      console.log("resp", resp);
    }, err => {
      console.warn("err", err);
    });
  }

  panRight = (cam: Camera) => {
    this.tilting = true
    this.timeout = 0
    console.log("panning right", cam.panRight);
    if (!cam.panRight) {
      return;
    }

    this.http.get(cam.panRight).subscribe(resp => {
      console.log("resp", resp);
    }, err => {
      console.warn("err", err);
    });
  }

  panTiltStop = (cam: Camera) => {
    if (!this.tilting) {
      return
    }

    console.log("stopping pan", cam.panTiltStop);
    if (!cam.panTiltStop) {
      return;
    }

    this.http.get(cam.panTiltStop).subscribe(resp => {
      console.log("resp", resp);
    }, err => {
      console.warn("err", err);
    });
    this.tilting = false
  }

  zoomIn = (cam: Camera) => {
    this.zooming = true
    this.timeout = 0
    console.log("zooming in", cam.zoomIn);
    if (!cam.zoomIn) {
      return;
    }

    this.http.get(cam.zoomIn).subscribe(resp => {
      console.log("resp", resp);
    }, err => {
      console.warn("err", err);
    });
  }

  zoomOut = (cam: Camera) => {
    this.zooming = true
    this.timeout = 0
    console.log("zooming out", cam.zoomOut);
    if (!cam.zoomOut) {
      return;
    }

    this.http.get(cam.zoomOut).subscribe(resp => {
      console.log("resp", resp);
    }, err => {
      console.warn("err", err);
    });
  }

  zoomStop = (cam: Camera) => {
    if (!this.zooming) {
      return
    }
    console.log("stopping zoom", cam.zoomStop);
    if (!cam.zoomStop) {
      return;
    }

    this.http.get(cam.zoomStop).subscribe(resp => {
      console.log("resp", resp);
    }, err => {
      console.warn("err", err);
    });
    this.zooming = true
  }

  selectPreset = (preset: CameraPreset) => {
    this.timeout = 0
    console.log("selecting preset", preset.displayName, preset.setPreset);
    if (!preset.setPreset) {
      return;
    }

    this.http.get(preset.setPreset).subscribe(resp => {
      console.log("resp", resp);
    }, err => {
      console.warn("err", err);
    });
  }

  getStreamURL = (cam: Camera) => {
    if (this.timeout >= 60) {
      return ""
    }
    // if (window.innerWidth > 750 && window.innerHeight > 750) {
      // console.log(cam.stream)
      return cam.stream
    // }

    return ""
  } 
}


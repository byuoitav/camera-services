import {Component, HostListener, ViewChild, ElementRef, OnInit, OnDestroy, AfterViewInit, EventEmitter} from '@angular/core';
import {Router, ActivatedRoute} from "@angular/router";
import {HttpClient} from '@angular/common/http';

import {Camera, Preset} from '../../services/api.service';
import { MatDialog } from '@angular/material/dialog';
import { PresetsDialog } from 'src/app/dialogs/presets/presets.component';
import { CookieService } from 'ngx-cookie-service';
import {JwtHelperService} from '@auth0/angular-jwt';
import { ErrorDialog } from '../../dialogs/error/error.dialog';
import { RebootDialog } from 'src/app/dialogs/reboot/reboot.component';


function isCameras(obj: Camera[] | any): obj is Camera[] {
  const cams = obj as Camera[];
  if (!cams || !cams.length || cams.length === 0) {
    return false;
  }

  return cams[0].displayName !== undefined;
}

@Component({
  selector: 'app-camera-feed',
  templateUrl: './camera-feed.component.html',
  styleUrls: ['./camera-feed.component.scss']
})
export class CameraFeedComponent implements OnInit, OnDestroy {
  rowHeight = "4:1.75";

  admin = false;
  rebooting = false;

  timeout = 0;
  cameras: Camera[];

  tilting = false;
  zooming = false;

  reboot: EventEmitter<boolean> = new EventEmitter();

  @ViewChild('stream') img: ElementRef;

  constructor(
    private router: Router,
    private http: HttpClient,
    public route: ActivatedRoute,
    private dialog: MatDialog,
    private cookieService: CookieService,
  ) {
    this.route.params.subscribe(params => {
      if ("room" in params && typeof params.room === "string") {
        const split = params.room.split("-");
        if (split.length == 2) {
          document.title = `${split[0]} ${split[1]} Camera Control`;
        } else {
          document.title = `${params.room} Camera Control`;
        }
      }
    })

    this.route.data.subscribe(data => {
      if ("cameras" in data && isCameras(data.cameras)) {
        this.cameras = data.cameras;
      }
    })

    this.reboot.subscribe(rebooting => {
      if (rebooting == true) {
        this.rebooting = true
        setTimeout(() => {
          this.rebooting = false;
          this.timeout = 0;
          window.location.reload();
        }, 45000);
      }
    })
  }

  ngOnInit() {
    if (window.innerWidth <= 1024) {
      this.rowHeight = "4:1.25";
    }

    const decoder = new JwtHelperService();
    var decoded = decoder.decodeToken(this.cookieService.get("camera-services-control"))
    if (decoded != null && decoded.auth.restart == true) {
        this.admin = true;
    }

    setInterval(() => {
      this.timeout++
      if (this.timeout == 60) {
        console.log("preview timing out")
      }
    }, 1000);
  }

  ngOnDestroy() {
    this.img.nativeElement.src = "";
  }

  @HostListener("window:resize", ["$event"])
  onResize(event) {
    if (window.innerWidth > 1024) {
      this.rowHeight = "4:1.75"
    } else {
      this.rowHeight = "4:1.25"
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
    this.zooming = false
  }

  selectPreset = (preset: Preset) => {
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
    if (this.timeout >= 60 || this.rebooting) {
      return ""
    }

    return cam.stream
  }

  openPresetsDialog = (cam: Camera) => {
    const dialogs = this.dialog.openDialogs.filter(dialog => {
      return dialog.componentInstance instanceof PresetsDialog
    })

    if (dialogs.length > 0) {
      return
    }

    this.dialog.open(PresetsDialog, {
      width: "fit-content",
      data: {
        presets: cam.presets
      }
    })
  }

  openRebootDialog = (cam: Camera) => {
    const dialogs = this.dialog.openDialogs.filter(dialog => {
      return dialog.componentInstance instanceof RebootDialog
    })

    if (dialogs.length > 0) {
      return
    }

    this.dialog.open(RebootDialog, {
      width: "fit-content",
      data: {
        camera: cam,
        reboot: this.reboot,
      }
    })  
  }

  checkSavePreset = (cam: Camera) => {
    return cam.presets.some(function (element) {
      return (element.savePreset != undefined)
    })
  }
}

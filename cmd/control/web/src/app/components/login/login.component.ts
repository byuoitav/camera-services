import {Component, OnInit, EventEmitter, ViewChild, ElementRef, AfterViewInit} from '@angular/core';
import {Router, ActivatedRoute} from '@angular/router';
import {HttpErrorResponse} from "@angular/common/http";
import {MatDialog} from "@angular/material/dialog";
import {APIService, Camera} from "../../services/api.service";

import {ErrorDialog} from "../../dialogs/error/error.dialog";
import { CookieService } from 'ngx-cookie-service';
import {JwtHelperService} from '@auth0/angular-jwt';
import {MatSnackBar} from '@angular/material/snack-bar';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit, AfterViewInit {
  cameras: Camera[]

  key = "";
  keyboardEmitter: EventEmitter<string>;

  @ViewChild('form') form: ElementRef;

  constructor(
    private router: Router,
    private api: APIService,
    private dialog: MatDialog,
    private route: ActivatedRoute,
    private cookieService: CookieService,
    private snackBar: MatSnackBar,
    ) {}

  ngOnInit() {
    this.keyboardEmitter = new EventEmitter<string>();
    this.keyboardEmitter.subscribe(s => {
      if (s === "done") {
        this.goToCameraControl()
      } else {
        this.key = s;
      }
    });
    document.title = "BYU Camera Control";
  }

  ngAfterViewInit() {
    this.form.nativeElement.focus();
    const decoder = new JwtHelperService();
    var decoded = decoder.decodeToken(this.cookieService.get("camera-services-control"))
    // console.log(decoded)
    // console.log("rooms ", Object.keys(this.rooms).length)
    // console.log("rooms", decoded.rooms)
    if (decoded != null && (Object.keys(decoded.rooms).length > 0)){
      var room = this.findMostRecentRoom(decoded.rooms)
      let snackBarRef = this.snackBar.open("Would you like to go back to " + room.name + "?", "GO", {duration: 10000,})
      snackBarRef.onAction().subscribe(() => {
        console.log("they pushed go!")
        this.router.navigate(["/control/" + room.name], {
          queryParams: {
            controlGroup: room.controlGroup
          }
        })
      })
    }
  }

  _focus() {
    this.form.nativeElement.focus();
  }

  hasFocus() {
    return document.hasFocus();
  }

  getCodeChar = (index: number): string => {
    if (this.key.length > index) {
      return this.key.charAt(index);
    }

    return "";
  }

  goToCameraControl = async () => {
    this.api.getControlInfo(this.key).subscribe(info => {
      return this.router.navigate(["/control/" + info.room], {
        queryParams: {
          controlGroup: info.controlGroup,
        }
      });
    }, (err: HttpErrorResponse) => {
      let msg = err.error;
      switch (err.status) {
        case 401:
          msg = "Invalid room code";
        default:
      }

      this.showError(msg);
    })

    this.key = "";
  }

  showError(msg: string) {
    const dialogs = this.dialog.openDialogs.filter(dialog => {
      return dialog.componentInstance instanceof ErrorDialog
    })

    if (dialogs.length > 0) {
      return;
    }

    this.dialog.open(ErrorDialog, {
      width: "fit-content",
      data: {
        msg: msg,
      }
    })
  }

  findMostRecentRoom(rooms: Map<string, Map<string, string>>) {
    var mostRecent: Room

    for (const [room, val] of Object.entries(rooms)) {
      console.log(room, val)
      for (const [preset, time] of Object.entries(val)) {
        console.log(preset, time)
        if (mostRecent === undefined || mostRecent.controlGroup.time < time) {
          mostRecent = new Room(room, new ControlGroup(preset, time as string))
        }
        console.log("most recent", mostRecent)
      }
    }

    return mostRecent
  }
}

export class Room {
  name: string
  controlGroup: ControlGroup
  constructor(name: string, controlGroup: ControlGroup) {
    this.name = name;
    this.controlGroup = controlGroup
  }
}

export class ControlGroup {
  name: string
  time: string

  constructor(name: string, time: string) {
    this.name = name;
    this.time = time;
  }
}

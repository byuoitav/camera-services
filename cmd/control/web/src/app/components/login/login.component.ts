import {Component, OnInit, EventEmitter, ViewChild, ElementRef, AfterViewInit} from '@angular/core';
import {Router, ActivatedRoute} from '@angular/router';
import {HttpErrorResponse} from "@angular/common/http";
import {MatDialog} from '@angular/material/dialog';
import {APIService, Camera} from "../../services/api.service";
import {MatFormFieldModule} from '@angular/material/form-field';
import {ErrorDialog} from "../../dialogs/error/error.dialog";
import { CookieService } from 'ngx-cookie-service';
import {JwtHelperService} from '@auth0/angular-jwt';
import {MatSnackBar as MatSnackBar} from '@angular/material/snack-bar';

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
    private matFormField: MatFormFieldModule,
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
    if (decoded != null && (Object.keys(decoded.rooms).length > 0)){
      var room = this.findMostRecentRoom(decoded.rooms)
      let snackBarRef = this.snackBar.open("Would you like to go back to " + room.name + "?", "GO", {duration: 30000,})
      snackBarRef.onAction().subscribe(() => {
        this.router.navigate(["/control/" + room.name], {
          queryParams: {
            controlGroup: room.controlGroup.name
          },
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
    this.cookieService.delete("control-key")
    this.cookieService.set("control-key", this.key, 1/24)

    this.api.getControlInfo(this.key).subscribe(info => {
      this.snackBar.dismiss();
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
      for (const [preset, time] of Object.entries(val)) {
        if (mostRecent === undefined || mostRecent.controlGroup.time < time) {
          mostRecent = new Room(room, new ControlGroup(preset, time as string))
        }
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

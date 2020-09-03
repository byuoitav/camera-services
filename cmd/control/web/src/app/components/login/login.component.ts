import {Component, OnInit, EventEmitter, ViewChild, ElementRef, AfterViewInit} from '@angular/core';
import {Router} from '@angular/router';
import {HttpErrorResponse} from "@angular/common/http";
import {MatDialog} from "@angular/material/dialog";
import {APIService, Camera} from "../../services/api.service";

import {ErrorDialog} from "../../dialogs/error/error.dialog";

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

  constructor(private router: Router, private api: APIService, private dialog: MatDialog) {}

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
}

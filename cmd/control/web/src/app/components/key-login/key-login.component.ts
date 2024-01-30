import {Component, OnInit, EventEmitter, ViewChild, ElementRef, AfterViewInit} from '@angular/core';
import {Router, ActivatedRoute} from '@angular/router';
import {HttpErrorResponse} from "@angular/common/http";
import { MatDialog } from '@angular/material/dialog';
import {APIService, Camera} from "../../services/api.service";

import {ErrorDialog} from "../../dialogs/error/error.dialog";

@Component({
  selector: 'key-login',
  templateUrl: './key-login.component.html',
  styleUrls: ['./key-login.component.scss']
})
export class KeyLoginComponent implements OnInit, AfterViewInit {
  key = "";

  @ViewChild('form') form: ElementRef;

  constructor(private router: Router, 
    private api: APIService, 
    private dialog: MatDialog, 
    private route: ActivatedRoute) {}

  ngOnInit() {
    document.title = "BYU Camera Control";
  }

  ngAfterViewInit() {

    // If a key was passed in, use it and immediately call goToCameraControl
    this.route.queryParams.subscribe(params => {
      if (params['key']) {
        this.key = params['key']
		    this.goToCameraControl()
      } else {
	      return this.router.navigate(["/login"])
      }
    });
  }

  _focus() {
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

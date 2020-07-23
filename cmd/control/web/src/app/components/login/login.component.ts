import {Component, OnInit, EventEmitter, ViewChild, ElementRef, ChangeDetectorRef, NgZone, AfterViewInit} from '@angular/core';
import {Router} from '@angular/router';
import {Camera} from '../../../objects/objects';
import { ɵINTERNAL_BROWSER_DYNAMIC_PLATFORM_PROVIDERS } from '@angular/platform-browser-dynamic';


@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit, AfterViewInit {
  cameras: Camera[]

  key = "";
  keyboardEmitter: EventEmitter<string>;

  @ViewChild('form') form:ElementRef;

  constructor(private router: Router) {}

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
    console.log("logging in with key", this.key);
    this.router.navigate(["/key/" + this.key]);
    this.key = "";
  }
}

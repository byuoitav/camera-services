import { Component, OnInit, EventEmitter } from '@angular/core';
import { MatBottomSheet } from '@angular/material/bottom-sheet';
import { Router, NavigationStart, NavigationError } from '@angular/router';
import { HttpClient } from "@angular/common/http";
import { Camera } from '../../../objects/objects';
import { NumpadDialog } from '../../dialogs/numpad/numpad.dialog';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {
  cameras: Camera[]

  key= "";
  keyboardEmitter: EventEmitter<string>;

  
  constructor(private router: Router, private http: HttpClient, private bottomSheet: MatBottomSheet){
    // this.router.events.subscribe(event => {
    //   switch (true) {
    //     case event instanceof NavigationStart:
    //       this.loggingIn = true;
    //       break;
    //     case event instanceof NavigationError:
    //       this.loggingIn = false;
    //       break;
    //   }
    // })
  }

  ngOnInit() {
    this.keyboardEmitter = new EventEmitter<string>();
    this.keyboardEmitter.subscribe(s => {
      if (s === "done") {
        this.goToCameraControl()
      } else {
        this.key = s;
      }
    });
  }

  codeKeyUp(event, index) {
    console.log(event);
    if (event.key === "Backspace") {
      if (index > 0) {
        const elementName = "codeKey" + (index + 1);
        document.getElementById(elementName).focus();
      }
      return;
    }
    if (index >= 0 && index < 5) {
      const elementName = "codeKey" + (index + 1);
      document.getElementById(elementName).focus(); 
    }
  }

  getCodeChar = (index: number): string => {
    if (this.key.length > index) {
      return this.key.charAt(index);
    }

    return "";
  }

  showNumpad() {
    console.log("trying to show numpad")
    this.bottomSheet.open(NumpadDialog, {
      data: this.keyboardEmitter,
    });
  }

  goToCameraControl = async () => {
    console.log("logging in with key", this.key);    
    this.router.navigate(["/key/" + this.key]);
    this.key = "";
  }
}

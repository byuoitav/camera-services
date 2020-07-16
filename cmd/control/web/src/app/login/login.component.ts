import { Component, OnInit, EventEmitter } from '@angular/core';
import { MatBottomSheet } from '@angular/material/bottom-sheet';
import { Router } from '@angular/router';
import { HttpClient } from "@angular/common/http";
import { Camera } from '../../objects/objects';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {
  cameras: Camera[]

  key= "";
  keyboardEmitter: EventEmitter<string>;
  private bottomSheet: MatBottomSheet

  constructor(private router: Router, private http: HttpClient){}

  ngOnInit() {
    this.keyboardEmitter = new EventEmitter<string>();
    this.keyboardEmitter.subscribe(s => {
      this.key = s;
    });
  }

  goToCameraControl = async () => {
    console.log("logging in with key", this.key);    
    this.router.navigate(["/key/" + this.key])
    this.key = "";
  }
}

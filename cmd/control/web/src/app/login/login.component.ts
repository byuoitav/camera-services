import { Component, OnInit, EventEmitter } from '@angular/core';
import { MatBottomSheet } from '@angular/material/bottom-sheet';
import { Router } from '@angular/router';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {

  key= "";
  keyboardEmitter: EventEmitter<string>;
  private bottomSheet: MatBottomSheet

  constructor(private router: Router){}

  ngOnInit() {
    this.keyboardEmitter = new EventEmitter<string>();
    this.keyboardEmitter.subscribe(s => {
      this.key = s;
    });
  }

  goToCameraControl = async () => {
    console.log("logging in with key", this.key);
    const success = await this.router.navigate(["/key/" + this.key])
    this.key = "";
  }
}

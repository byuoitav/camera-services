import { Component, HostListener, ViewChild, ElementRef, OnInit } from '@angular/core';
import { Router } from "@angular/router";

@Component({
  selector: 'app-camera-feed',
  templateUrl: './camera-feed.component.html',
  styleUrls: ['./camera-feed.component.scss']
})
export class CameraFeedComponent implements OnInit {
  rowHeight = "4:1.75"
  timeout = 0  
  constructor(private router: Router) {}
  ngOnInit() {
    setInterval(() => {
      this.timeout++
      if (this.timeout == 60) {
        console.log("preview timing out")
      }
    }, 1000)
  }


  @HostListener("window:resize", ["$event"])
  onResize(event) {
    if (window.innerWidth >= 1024 && window.innerHeight >= 768 && window.innerHeight <= 1024) {
      this.rowHeight = "4:2.5"
    } else {
      this.rowHeight = "4:1.75"
    }
  }

  exitRoom() {
    console.log("exiting room")
    this.router.navigate([""])
  }
}


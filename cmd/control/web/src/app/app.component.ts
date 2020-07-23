import {Component} from '@angular/core';
import {Router, NavigationStart, NavigationEnd, NavigationCancel} from '@angular/router';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {

  loading = false;

  constructor(private router: Router) {
    let vh = window.innerHeight * 0.01;

    this.router.events.subscribe(event => {
      if (event instanceof NavigationStart) {
        if (event.url.includes("/key/")) {
          this.loading = true;
        }
      }
      if (event instanceof NavigationEnd) {
        if (event.url.includes("/key/")) {
          this.loading = false;
        }
      }
      if (event instanceof NavigationCancel) {
        this.loading = false;
      }
    });

  }
}

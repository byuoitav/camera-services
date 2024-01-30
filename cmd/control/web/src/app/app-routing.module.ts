import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {AppComponent} from './app.component';
import {CameraFeedComponent} from './components/camera-feed/camera-feed.component';
import {LoginComponent} from './components/login/login.component';
import {KeyLoginComponent} from './components/key-login/key-login.component';
import {RoomResolver} from './services/room.resolver';


const routes: Routes = [
  {
    path: "",
    redirectTo: "/login",
    pathMatch: "full"
  },
  {
    path: "",
    component: AppComponent,
    children: [
      {
        path: "login",
        component: LoginComponent
      },
      {
        path: "key-login",
        component: KeyLoginComponent
      },
      {
        path: "control/:room",
        resolve: {
          cameras: RoomResolver
        },
        children: [
          {
            path: "",
            component: CameraFeedComponent
          }
        ]
      }
    ]
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes, {})],
  exports: [RouterModule]
})
export class AppRoutingModule {}

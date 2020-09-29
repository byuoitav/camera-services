import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {MatButtonModule} from "@angular/material/button";
import {MatTabsModule} from "@angular/material/tabs";
import {MatIconModule} from "@angular/material/icon";
import {MatGridListModule} from "@angular/material/grid-list";
import {MatDividerModule} from "@angular/material/divider";
import {CameraFeedComponent} from './components/camera-feed/camera-feed.component';
import {MatFormFieldModule} from '@angular/material/form-field';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {MatInputModule} from '@angular/material/input';
import {MatBottomSheetModule} from '@angular/material/bottom-sheet';
import {LoginComponent} from './components/login/login.component';
import {HttpClientModule} from '@angular/common/http';
import {ErrorDialog} from './dialogs/error/error.dialog';
import {MatDialogModule} from '@angular/material/dialog';
import {CookieService} from 'ngx-cookie-service';
import {MatSnackBarModule} from '@angular/material/snack-bar';
import {PresetsDialog} from './dialogs/presets/presets.component';
import { RebootDialog } from './dialogs/reboot/reboot.component';

@NgModule({
  declarations: [
    AppComponent,
    CameraFeedComponent,
    LoginComponent,
    ErrorDialog,
    PresetsDialog,
    RebootDialog,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    MatTabsModule,
    MatButtonModule,
    MatIconModule,
    MatGridListModule,
    MatDividerModule,
    MatFormFieldModule,
    FormsModule,
    ReactiveFormsModule,
    MatInputModule,
    MatBottomSheetModule,
    HttpClientModule,
    MatDialogModule,
    MatSnackBarModule,
  ],
  providers: [CookieService],
  bootstrap: [AppComponent],
  entryComponents: [
    ErrorDialog,
  ]
})
export class AppModule {}

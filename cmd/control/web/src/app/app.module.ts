import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';
import {Router, ActivatedRoute} from "@angular/router";
import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
//import {MatLegacyButtonModule as MatButtonModule} from "@angular/material/legacy-button";
import {MatButtonModule as MatButtonModule} from "@angular/material/button";

import {MatIconModule} from "@angular/material/icon";
import {MatGridListModule} from "@angular/material/grid-list";
import {MatDividerModule} from "@angular/material/divider";
import {CameraFeedComponent} from './components/camera-feed/camera-feed.component';

import {FormsModule, ReactiveFormsModule} from '@angular/forms';
//import {MatLegacyInputModule as MatInputModule} from '@angular/material/legacy-input';
import {MatInputModule} from '@angular/material/input';
import {MatBottomSheetModule} from '@angular/material/bottom-sheet';
import {LoginComponent} from './components/login/login.component';
import {HttpClientModule} from '@angular/common/http';
import {ErrorDialog} from './dialogs/error/error.dialog';
import {MatDialogModule} from '@angular/material/dialog';
import {CookieService} from 'ngx-cookie-service';
//import {MatLegacySnackBarModule as MatSnackBarModule} from '@angular/material/legacy-snack-bar';
import {MatSnackBarModule} from '@angular/material/snack-bar';
import {PresetsDialog} from './dialogs/presets/presets.component';
import {RebootDialog} from './dialogs/reboot/reboot.component';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatTabsModule} from '@angular/material/tabs';



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
        MatFormFieldModule,
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
    bootstrap: [AppComponent]
})
export class AppModule {}

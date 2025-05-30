import { Component, OnInit, Inject, Output, EventEmitter } from '@angular/core';
import { Camera } from 'src/app/services/api.service';
//import { MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA, MatLegacyDialogRef as MatDialogRef, MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';
import { MAT_DIALOG_DATA } from '@angular/material/dialog';
import { MatDialogRef } from '@angular/material/dialog';
import { MatDialog } from '@angular/material/dialog';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-reboot',
  templateUrl: './reboot.component.html',
  styleUrls: ['./reboot.component.scss']
})
export class RebootDialog {


  constructor(
    private http: HttpClient,
    private dialog: MatDialog,
    public ref: MatDialogRef<RebootDialog>,
    @Inject(MAT_DIALOG_DATA)
    public data: {
      camera: Camera;
      reboot: EventEmitter<boolean>;
    }
  ) {}

  close = () => {
    this.ref.close();
  }

  confirm = () => {
    console.log("rebooting...")
    this.http.get(this.data.camera.reboot).subscribe()
    this.data.reboot.emit(true)
    this.ref.close();
  }

}

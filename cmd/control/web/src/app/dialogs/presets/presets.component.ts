import { Component, OnInit, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialog} from '@angular/material/dialog';
import { Preset } from 'src/app/services/api.service';
import { HttpClient } from '@angular/common/http';
import { ErrorDialog } from '../error/error.dialog';



@Component({
  selector: 'app-presets',
  templateUrl: './presets.component.html',
  styleUrls: ['./presets.component.scss']
})
export class PresetsDialog {
  curPreset: Preset;

  constructor(
    private http: HttpClient,
    private dialog: MatDialog,
    public ref: MatDialogRef<PresetsDialog>,
    @Inject(MAT_DIALOG_DATA)
    public data: {
      presets: Preset[];
    }
  ) {}

  close = () => {
    this.ref.close();
  }

  confirm = () => {
    if (this.curPreset == undefined) {
      return
    }
    console.log(this.curPreset)
    this.http.get(this.curPreset.setPreset).subscribe(resp => {
      console.log("resp", resp);
      this.ref.close();
    }, err => {
      console.warn("err", err);
      this.dialog.open(ErrorDialog, {
        data: {
          msg: "Unable to set preset"
        }
      })
    });
  }

  setCurPreset = (preset: Preset) => {
    var selected = document.querySelectorAll(".selected");
    for (let i = 0; i < selected.length; i++) {
      selected[i].classList.remove("selected");
    }
    document.getElementById(preset.displayName).classList.add("selected");
    this.curPreset = preset;
    console.log(this.curPreset);
  }

  disabled = () => {
    if (this.curPreset == undefined) {
      return true;
    }

    return false;
  }
}

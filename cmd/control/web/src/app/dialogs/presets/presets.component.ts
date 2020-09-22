import { Component, OnInit, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA} from '@angular/material/dialog';
import { Preset } from 'src/app/services/api.service';


@Component({
  selector: 'app-presets',
  templateUrl: './presets.component.html',
  styleUrls: ['./presets.component.scss']
})
export class PresetsDialog implements OnInit {
  curPreset: Preset;

  constructor(
    public ref: MatDialogRef<PresetsDialog>,
    @Inject(MAT_DIALOG_DATA)
    public data: {
      presets: Preset[];
    }
  ) {
    console.log(this.data.presets)
   }

  ngOnInit(): void {
  }

  close = () => {
    this.ref.close();
  }

  confirm = () => {
    if (this.curPreset == undefined) {
      return
    }
    console.log(this.curPreset)
    for (let i = 0; i < this.data.presets.length; i++) {
      if (this.curPreset == this.data.presets[i]) {
        //hit the endpoint here with i as the preset
      }
    }
  }

  setCurPreset = (preset: Preset) => {
    var selected = document.querySelectorAll(".selected");
    for (let i = 0; i < selected.length; i++) {
      selected[i].classList.remove("selected");
    }
    document.getElementById(preset.displayName).classList.add("selected");
    this.curPreset = preset;
  }

}

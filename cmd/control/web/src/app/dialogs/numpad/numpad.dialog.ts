import { Component, OnInit, AfterViewInit, Inject, EventEmitter, ViewEncapsulation } from '@angular/core';
import Keyboard from 'simple-keyboard';
import { MatBottomSheetRef, MAT_BOTTOM_SHEET_DATA } from '@angular/material/bottom-sheet'

@Component({
  selector: 'numpad',
  encapsulation: ViewEncapsulation.None,
  templateUrl: './numpad.dialog.html',
  styleUrls: [
    './numpad.dialog.scss',
    '../../../../node_modules/simple-keyboard/build/css/index.css']
})
export class NumpadDialog implements OnInit, AfterViewInit {

  private keyboard: Keyboard;
  roomCodeValue = '';

  constructor(
    private bottomSheetRef: MatBottomSheetRef<NumpadDialog>,
    @Inject(MAT_BOTTOM_SHEET_DATA) public data: EventEmitter<string>) {
   }

  ngOnInit() {
  }

  ngAfterViewInit() {
    this.keyboard = new Keyboard({
      onChange: input => this.onChange(input),
      onKeyPress: button => this.onKeyPress(button),
      layout: {
        default: [
          '1 2 3',
          '4 5 6',
          '7 8 9',
          '{bksp} 0 {enter}'
        ]
      },
      display: {
        '{bksp}': '⌫',
        '{enter}': 'OK'
      },
      maxLength: {
        default: 6
      }
    });

    this.keyboard.addButtonTheme('{enter}', 'kb-done');
  }

  onChange = (input: string) => {
    this.roomCodeValue = input;
    this.data.emit(this.roomCodeValue);
  }

  onKeyPress = (button: string) => {
    // this.roomCode.nativeElement.focus();
    if (button === '{bksp}') {
      this.roomCodeValue = this.roomCodeValue.substring(0, this.roomCodeValue.length - 1);
    }
    if (button === '{enter}') {
      this.data.emit("done")
      this.bottomSheetRef.dismiss('all good in the hood');
    }
  }

}

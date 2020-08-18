import {Component, Injectable, OnInit} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {FormGroup, FormBuilder, Validators} from "@angular/forms";
import {Observable, of} from "rxjs";
import {tap, startWith, debounceTime, distinctUntilChanged, switchMap, map} from "rxjs/operators";
import {StepperSelectionEvent} from "@angular/cdk/stepper";

@Injectable({
  providedIn: 'root',
})
export class APIService {
  constructor(private http: HttpClient) {}

  rooms = [];
  presets = [];

  getRooms() {
    return this.rooms.length ?
      of(this.rooms) :
      this.http.get<string[]>("/api/v1/rooms").pipe(tap(data => this.rooms = data));
  }

  getControlGroups(room: string) {
    return this.http.get<string[]>("/api/v1/rooms/" + encodeURIComponent(room) + "/controlGroups");
  }
}

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  firstFormGroup: FormGroup;
  secondFormGroup: FormGroup;
  filteredRooms: Observable<any[]>;
  controlGroups: string[];

  constructor(private _api: APIService, private _formBuilder: FormBuilder) {
    this.firstFormGroup = this._formBuilder.group({
      room: ['', Validators.required],
    })

    this.secondFormGroup = this._formBuilder.group({
      controlGroup: ['', Validators.required],
    })
  }

  ngOnInit() {
    this.filteredRooms = this.firstFormGroup.controls['room'].valueChanges.pipe(
      startWith(''),
      debounceTime(250),
      distinctUntilChanged(),
      switchMap(val => {
        return this.filter(val || '');
      })
    )
  }

  filter(val: string): Observable<any[]> {
    return this._api.getRooms()
      .pipe(
        map(resp => resp.filter(opt => {
          return opt.toLowerCase().indexOf(val.toLowerCase()) === 0;
        }))
      );
  }

  formStep(event: StepperSelectionEvent) {
    switch (event.selectedIndex) {
      case 0:
        this.controlGroups = [];
      case 1:
        const room = this.firstFormGroup.controls['room'].value;
        this._api.getControlGroups(room).subscribe(data => {
          this.controlGroups = data;
        })
    }
  }

  selectControlGroup(cg: string) {
    let room = this.firstFormGroup.controls['room'].value;
    room = encodeURIComponent(room)
    cg = encodeURIComponent(cg)

    window.location.assign('/api/v1/rooms/' + room + '/controlGroups/' + cg);
  }
}

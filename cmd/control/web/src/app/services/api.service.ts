import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";

export interface ControlInfo {
  room: string;
  controlGroup: string;
  controlKey: string;
}

export interface Camera {
  displayName: string;

  tiltUp: string;
  tiltDown: string;
  panLeft: string;
  panRight: string;
  panTiltStop: string;

  zoomIn: string;
  zoomOut: string;
  zoomStop: string;

  reboot: string;

  memoryRecall: string;
  stream: string;

  presets: Preset[];
}

export interface Preset {
  displayName: string;
  setPreset: string;
  savePreset: string;
}

@Injectable({
  providedIn: 'root'
})
export class APIService {
  constructor(private http: HttpClient) {}

  getControlInfo(key: string) {
    return this.http.get<ControlInfo>("/api/v1/controlInfo", {
      params: {
        key: key,
      }
    })
  }

  getCameras(room: string, controlKey: string, controlGroup: string) {
    return this.http.get<Camera[]>("/api/v1/cameras", {
      params: {
        room: room,
        controlGroup: controlGroup,
        controlKey: controlKey,
      }
    })
  }
}

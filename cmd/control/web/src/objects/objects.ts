export class Camera {
    displayName: string;

    tiltUp: string;
    tiltDown: string;
    panLeft: string;
    panRight: string;
    panTiltStop: string;
  
    zoomIn: string;
    zoomOut: string;
    zoomStop: string;
  
    memoryRecall: string;
  
    presets: CameraPreset[];
}

export class CameraPreset {
    displayName: string;
    setPreset: string;
}

export class Config {
    cameras: Camera[]
}
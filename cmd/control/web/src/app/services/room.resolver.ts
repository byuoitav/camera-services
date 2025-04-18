import {Injectable} from "@angular/core";
import { Router, ActivatedRouteSnapshot, RouterStateSnapshot } from "@angular/router";
import {Observable, of, EMPTY} from "rxjs";
import {HttpErrorResponse} from "@angular/common/http";
import {take, mergeMap, catchError} from 'rxjs/operators';
import { MatDialog } from '@angular/material/dialog';
import { CookieService } from 'ngx-cookie-service';


import {APIService, Camera} from "./api.service";
import {ErrorDialog} from "../dialogs/error/error.dialog";

@Injectable({
  providedIn: "root"
})
export class RoomResolver  {
  constructor(
    private router: Router, 
    private api: APIService, 
    private dialog: MatDialog,
    private cookieService: CookieService,
  ) {}

  resolve(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot
  ): Observable<Camera[]> | Observable<never> {
    const room = route.paramMap.get("room");
    const cg = route.queryParamMap.get("controlGroup");
    const controlKey = this.cookieService.get("control-key");

    return this.api.getCameras(room, controlKey, cg).pipe(
      take(1),
      catchError((err: HttpErrorResponse) => {
        this.router.navigate([""], {});

        let msg = err.error;
        switch (err.status) {
          case 401:
            msg = `Not authorized to control ${room}`;
          default:
        }

        this.showError(msg);
        return EMPTY;
      }),
      mergeMap(cameras => {
        return of(cameras);
      })
    );
  }

  showError(msg: string) {
    const dialogs = this.dialog.openDialogs.filter(dialog => {
      return dialog.componentInstance instanceof ErrorDialog
    })

    if (dialogs.length > 0) {
      return;
    }

    this.dialog.open(ErrorDialog, {
      width: "fit-content",
      data: {
        msg: msg,
      }
    })
  }
}

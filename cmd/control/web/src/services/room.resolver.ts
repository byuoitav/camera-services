import {Injectable, EventEmitter} from "@angular/core";
import {Router, Resolve, ActivatedRouteSnapshot, RouterStateSnapshot} from "@angular/router";
import {Config} from "../objects/objects";
import {Observable, of} from "rxjs";
import {HttpClient, HttpRequest, HttpHeaders} from "@angular/common/http";
import {catchError, tap, map} from 'rxjs/operators';
import {MatDialog} from "@angular/material/dialog";
import {ErrorDialog} from "../app/dialogs/error/error.dialog";


@Injectable({
  providedIn: "root"
})
export class RoomResolver implements Resolve<Config>{
  constructor(private router: Router, private http: HttpClient, private dialog: MatDialog) {}

  resolve(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot
  ): Observable<Config> | Observable<never> {
    const key = route.paramMap.get("key");
    return this.http.get<Config>("/api/v1/key/" + key).pipe(
      catchError(err => {
        console.log("error", err)
        this.router.navigate([""], {})

        const dialogs = this.dialog.openDialogs.filter(dialog => {
          return dialog.componentInstance instanceof ErrorDialog
        })

        if (dialogs.length < 1) {
          const ref = this.dialog.open(ErrorDialog, {
            width: "fit-content",
            data: {
              msg: "No cameras found for given room key"
            }
          })
        }

        return Observable.throw(err)
      }),
      tap(data => {
        console.log("resolver data", data)
      }),
    )
  }
}

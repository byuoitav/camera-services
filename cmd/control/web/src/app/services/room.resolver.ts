import {Injectable} from "@angular/core";
import {Router, Resolve, ActivatedRouteSnapshot, RouterStateSnapshot} from "@angular/router";
import {Observable, of, EMPTY} from "rxjs";
import {take, mergeMap, catchError} from 'rxjs/operators';

import {APIService, Camera} from "./api.service";


@Injectable({
  providedIn: "root"
})
export class RoomResolver implements Resolve<Camera[]> {
  constructor(private router: Router, private api: APIService) {}

  resolve(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot
  ): Observable<Camera[]> | Observable<never> {
    const room = route.paramMap.get("room");
    const cg = route.queryParamMap.get("controlGroup");

    return this.api.getCameras(room, cg).pipe(
      take(1),
      mergeMap(cameras => {
        return of(cameras);
      })
    );
  }
}

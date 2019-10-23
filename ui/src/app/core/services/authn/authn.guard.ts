import {Injectable} from '@angular/core';
import {ActivatedRouteSnapshot, CanActivate, Router, RouterStateSnapshot} from '@angular/router';
import {Observable} from 'rxjs';
import {AuthnService} from "./authn.service";
import {map, take} from "rxjs/operators";

@Injectable()
export class AuthnGuard implements CanActivate {

  constructor(
    private authnService: AuthnService,
    private router: Router
  ) {}

  canActivate(
    next: ActivatedRouteSnapshot,
    state: RouterStateSnapshot): Observable<boolean> | Promise<boolean> | boolean {
    return this.authnService.isLoggedIn
      .pipe(
        take(1),
        map((isLoggedIn: boolean) => {
          console.log(`logged-in state: ${isLoggedIn}`);
          if (isLoggedIn) {
            return true;
          }
          this.router.navigate(['/login']);
          return false;
        })
      );
  }
}

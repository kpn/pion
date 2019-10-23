import {Injectable, isDevMode} from '@angular/core';
import {BehaviorSubject, of, throwError} from "rxjs";
import {Router} from "@angular/router";
import {User} from "./user";
import {HttpClient} from "@angular/common/http";
import {tap} from "rxjs/operators";
import {environment} from "../../../../environments/environment";

@Injectable()
export class AuthnService {
  private loggedIn = new BehaviorSubject<boolean>(false);
  private loginUrl = '/public/login';
  private logoutUrl = '/restricted/logout';

  private loggedInUser;

  get isLoggedIn() {
    return this.loggedIn.asObservable();
  }

  constructor(private http: HttpClient,
              private router: Router) {
  }

  login(user: User) {
    if (isDevMode() && !environment.remote) {
      // only in dev mode
      if (user.username === 'a' && user.password === 'a') {
        this.loggedIn.next(true);
        this.router.navigateByUrl('/main');
        const demoUser = {
          username: 'ngo500',
          firstName: 'Canh',
          lastName: 'Ngo',
          displayName: 'Ngo, Canh',
          title: 'Engineer',
          mail: 'canh.ngo@kpn.com',
          customer: 'Infraplatform',
          userGroups: ["dig_infraplatform"]
        };
        this.loggedInUser = demoUser;
        return of(demoUser);
      } else {
        return throwError('login failed')
      }
    }

    return this.http.post(this.loginUrl, user)
      .pipe(
        tap(loggedInUser => {
          this.loggedInUser = loggedInUser;
            this.log('login succeeded:', loggedInUser);
            this.loggedIn.next(true);
            this.router.navigateByUrl('/main');
          },
          err => this.error('login failed', err))
      );
  }

  logout() {
    return this.http.post(this.logoutUrl, null)
      .pipe(tap(() => {
          this.loggedIn.next(false);
          this.router.navigateByUrl('/login');
        },
        err => this.error('Logout failed', err) // TODO error handling
      ));
  }

  get currentUser() {
    return this.loggedInUser
  }

  private log(message: string, obj: Object) {
    console.log(message);
    console.log(obj)
  }

  private error(message: string, obj: Object) {
    console.error(message);
    console.error(obj)
  }
}

import {Injectable} from '@angular/core';
import {Observable, of} from "rxjs";
import {catchError, tap} from "rxjs/operators";
import {HttpClient} from "@angular/common/http";
import {Role} from "../model/role";

@Injectable({
  providedIn: 'root'
})
export class RoleService {
  private serviceUrl = 'restricted/roles';

  constructor(private http: HttpClient) {
  }

  getRoles(): Observable<Role[]> {
    return this.http.get<Role[]>(this.serviceUrl)
      .pipe(
        tap(roles => this.log(`fetched roles: ${roles}`)),
        catchError(this.handleError('getRoles', []))
      );
  }

  getRole(name: string): Observable<Role> {
    return this.http.get<Role>(`${this.serviceUrl}/${name}`)
      .pipe(
        tap(_ => this.log(`fetched role name=${name}`))
      );
  }

  private handleError<T>(operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {
      console.error(error); // log to console instead
      this.log(`${operation} failed: ${error.message}`);
      return of(result as T);
    };
  }

  private log(msg: string) {
    console.log(msg);
  }
}

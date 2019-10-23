import {Injectable} from '@angular/core';
import {Observable, of, throwError} from "rxjs";
import {RoleBinding} from "../model/role-binding";
import {catchError, tap} from "rxjs/operators";
import {HttpClient} from "@angular/common/http";

@Injectable({
  providedIn: 'root'
})
export class RoleBindingService {
  private serviceUrl = 'restricted/role_bindings';

  constructor(private http: HttpClient) {
  }

  getRoleBindings(): Observable<RoleBinding[]> {
    return this.http.get<RoleBinding[]>(this.serviceUrl)
      .pipe(
        tap(keys => this.log('fetched access keys')),
        catchError(this.handleError('getAccessKeys', []))
      );
  }

  createRoleBinding(roleBinding: RoleBinding): Observable<RoleBinding> {
    return this.http.post<RoleBinding>(this.serviceUrl, roleBinding)
      .pipe(
        tap(rb => this.log(`created new role-binding: ${rb.name}`))
      );
  }

  private log(msg: string) {
    console.log(msg);
  }

  deleteRoleBinding(name: string) {
    return this.http.delete(`${this.serviceUrl}/${name}`)
      .pipe(
        tap(_ => this.log(`deleted role-binding: ${name}`))
      );
  }

  private handleError<T>(operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {
      console.error(error); // log to console instead
      this.log(`${operation} failed: ${error.message}`);
      return of(result as T);
    };
  }

  updateRoleBinding(roleBinding: RoleBinding): Observable<RoleBinding> {
    return throwError('not implemented');
  }
}

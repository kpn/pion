import {Injectable} from '@angular/core';
import {AccessKey} from '../model/accesskey';
import {Observable, of, throwError} from 'rxjs';
import {HttpClient} from '@angular/common/http';
import {catchError, tap} from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class AccessKeyService {
  private serviceUrl = 'restricted/access_keys';  // URL to web api

  constructor(
    private http: HttpClient,
  ) {
  }

  deleteKey(key: AccessKey): Observable<AccessKey> {
    return this.http.delete<AccessKey>(`${this.serviceUrl}/${key.accessKeyId}`)
      .pipe(
        tap(_ => this.log(`deleted access key ${key.accessKeyId}`)),
        catchError(this.handleError<AccessKey>('deleteKey'))
      );
  }

  getAccessKeys(): Observable<AccessKey[]> {
    return this.http.get<AccessKey[]>(this.serviceUrl)
      .pipe(
        tap(keys => this.log('fetched access keys')),
        catchError(this.handleError('getAccessKeys', []))
      );
  }

  createKey(): Observable<AccessKey> {
    return this.http.post<AccessKey>(this.serviceUrl, {})
      .pipe(
        tap(key => this.log(`created new key: ${key.accessKeyId}`)),
        catchError(this.handleError('createKey', null))
      );
  }

  private log(msg: string) {
    console.log(msg);
  }

  private handleError<T>(operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {
      console.error(error); // log to console instead
      this.log(`${operation} failed:`);
      this.log(error);
      return throwError(result as T);
    };
  }

}

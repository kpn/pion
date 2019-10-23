import {Injectable} from '@angular/core';
import {Bucket} from '../model/bucket';
import {Observable, of} from 'rxjs';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import {catchError, tap} from 'rxjs/operators';

const headers = new HttpHeaders().set('Content-Type', 'application/json');

@Injectable({
  providedIn: 'root'
})
export class BucketService {
  private serviceUrl = 'restricted/buckets';

  constructor(private http: HttpClient) {
  }

  getBuckets(): Observable<Bucket[]> {
    console.log('calling bucket service');
    return this.http.get<Bucket[]>(this.serviceUrl)
      .pipe(
        tap(_ => this.log('fetched buckets')),
        catchError(this.handleError('getBuckets', []))
      );
  }

  createBucket(bucketName: string) {
    let body = {
      name: bucketName,
    };
    return this.http.post<Bucket>(this.serviceUrl, body, {headers: headers})
      .pipe(tap(()=>this.log(`created bucket ${bucketName}`)));
  }

  updateAcls(bucketName: string, acls: string): Observable<Bucket> {
    console.log(`Updating ACLs of bucket ${bucketName}: ${acls}`);
    return this.http.put<Bucket>(`${this.serviceUrl}/${bucketName}/acl`, acls,
      {headers: headers})
      .pipe(
        tap(_ => this.log(`updated ACLs for bucket ${bucketName}`)),
        catchError(this.handleError('updateAcls', null))
      );
  }

  private log(msg: string) {
    console.log(msg);
  }

  private handleError<T>(operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {
      console.error(error);
      this.log(`${operation} failed: status=${error.status}`);

      return of(result as T);
    };
  }
}

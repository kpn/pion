import {Component, OnInit, ViewChild} from '@angular/core';
import {BucketService} from '../core/services/bucket.service';
import {Bucket} from '../core/model/bucket';
import {MatDialog, MatSort, MatTableDataSource} from '@angular/material';
import {AclEditorDialogComponent} from '../acl-editor-dialog/acl-editor-dialog.component';
import {NewBucketDialogComponent} from "../new-bucket-dialog/new-bucket-dialog.component";

@Component({
  selector: 'app-buckets',
  templateUrl: './buckets.component.html',
  styleUrls: ['./buckets.component.scss']
})
export class BucketsComponent implements OnInit {
  displayedColumns: string[] = ['name', 'createdAt', 'modifiedAt', 'creator', 'acls'];
  dsBuckets: MatTableDataSource<Bucket>;

  @ViewChild(MatSort) sort: MatSort;

  constructor(
    private bucketService: BucketService,
    private dialog: MatDialog) {
  }

  ngOnInit() {
    this.getBuckets();
  }

  getBuckets(): void {
    this.bucketService.getBuckets()
      .subscribe(
        buckets => this.updateDataSource(buckets),
        err => console.log(`fetching buckets failed: ${err}`));
  }

  newBucket() {
    const dialogRef = this.dialog.open(NewBucketDialogComponent, {
      width: '500px',
      height: '200px'
    });

    dialogRef.afterClosed().subscribe(result => {
      if (!result) {
        return;
      }
      console.log(`Creating bucket ${result}`);
      this.bucketService.createBucket(result)
        .subscribe(
          bkt => this.addNewBucket(bkt),
          errResp => this.onAddBucketError(errResp)
        )
    });
  }

  private updateDataSource(buckets: Bucket[]) {
    this.dsBuckets = new MatTableDataSource(buckets ? buckets : []);
    this.dsBuckets.sort = this.sort
  }

  openAclDialog(bucket: Bucket): void {
    let aclsStr = bucket.acls ? JSON.stringify(bucket.acls) : "";
    const dialogRef = this.dialog.open(AclEditorDialogComponent, {
      width: '600px',
      height: '400px',
      data: {bucketName: bucket.name, acls: aclsStr}
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result && result.action === 'save') {
        this.saveACLs(bucket.name, result.data);
      }
    });
  }

  private saveACLs(bucketName: string, acls: string) {
    this.bucketService.updateAcls(bucketName, acls).subscribe(
      updatedBucket => this.updateBucket(bucketName, updatedBucket)
    );
  }

  private updateBucket(bucketName: string, updatedBucket: Bucket) {
    let targetBucket = this.dsBuckets.data.find(b => b.name === bucketName);
    // update object attributes without loosing reference
    Object.assign(targetBucket, updatedBucket);
  }

  private addNewBucket(bucket: Bucket) {
    this.dsBuckets.data.push(bucket);
    this.dsBuckets.data = this.dsBuckets.data.slice();
  }

  private onAddBucketError(errResp) {
    if (errResp == null) {
      console.log('Unknown error');
      return
    }
    let err = errResp.error;
    let errMessage;
    if (err && err.error_code) {
      switch (err.error_code) {
        case 'BucketAlreadyExists': {
          errMessage = 'Bucket already exist';
          break;
        }
        case 'BucketAlreadyOwnedByYou': {
          errMessage = 'Bucket already owned by you';
          break;
        }
        case 'InternalError': {
          errMessage = 'Internal error. Please report problem';
          break;
        }
        default: {
          // you should not go here
          errMessage = 'Unknown error. Please report the problem';
          break;
        }
      }
    } else {
      errMessage = `Unknown error: ${err}`
    }
    alert(`Creating new bucket failed: ${errMessage}`);
    console.log(err)
  }
}

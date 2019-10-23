import {Component, OnInit, ViewChild} from '@angular/core';
import {AccessKey} from '../core/model/accesskey';
import {AccessKeyService} from '../core/services/access-key.service';
import {MatDialog, MatSort, MatTableDataSource} from "@angular/material";
import {NewAccessKeyDialogComponent} from "../new-access-key-dialog/new-access-key-dialog.component";

@Component({
  selector: 'app-accesskeys',
  templateUrl: './access-keys.component.html',
  styleUrls: ['./access-keys.component.scss']
})
export class AccessKeysComponent implements OnInit {
  displayedColumns: string[] = ['createdAt', 'accessKeyId', 'lastUsed', 'status', 'actions'];
  dsAccessKeys: MatTableDataSource<AccessKey>;

  constructor(
    private accessKeyService: AccessKeyService,
    private dialog: MatDialog) {
  }

  @ViewChild(MatSort) sort: MatSort;

  ngOnInit() {
    this.getAccessKeys();
  }

  getAccessKeys(): void {
    this.accessKeyService.getAccessKeys()
      .subscribe(keys => this.updateDataSource(keys),
        err => console.log(`fetching access-keys failed: ${err}`));
  }

  private updateDataSource(keys) {
    this.dsAccessKeys = new MatTableDataSource(keys ? keys : []);
    this.dsAccessKeys.sort = this.sort;
  }

  deleteKey(key: AccessKey): void {
    if (confirm(`Are you sure to delete the access key '${key.accessKeyId}'?`)) {
      this.accessKeyService.deleteKey(key)
        .subscribe(() => this.dsAccessKeys.data = this.dsAccessKeys.data.filter(k => k !== key));
    }
  }

  newAccessKey(): void {
    this.accessKeyService.createKey()
      .subscribe(
        key => {
          this.showNewAccessKeyDialog(key);
          this.dsAccessKeys.data.push(key);
          this.dsAccessKeys.data = this.dsAccessKeys.data.slice();
        },
        err => console.error(`failed to create access key: ${err}`)
      );
  }

  private showNewAccessKeyDialog(key: AccessKey) {
    this.dialog.open(NewAccessKeyDialogComponent, {
      width: '600px',
      data: {accessKeyId: key.accessKeyId, secretKeyId: key.secretKeyId}
    });
  }
}

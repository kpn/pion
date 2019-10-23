import {Component, Inject} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material';

export interface AclDialogData {
  bucketName: string;
  acls: string;
}

@Component({
  selector: 'app-acl-editor-dialog',
  templateUrl: './acl-editor-dialog.component.html',
  styleUrls: ['./acl-editor-dialog.component.scss']
})
export class AclEditorDialogComponent {
  constructor(
    public dialogRef: MatDialogRef<AclEditorDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: AclDialogData) {
  }

  close(): void {
    this.dialogRef.close();
  }

  save(): void {
    this.dialogRef.close({action: 'save', data: this.data.acls});
  }
}

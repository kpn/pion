import {Component, Inject, OnInit} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material";

export interface AccessKeyData {
  accessKeyId: string
  secretKeyId: string
}

@Component({
  selector: 'app-new-access-key-dialog',
  templateUrl: './new-access-key-dialog.component.html',
  styleUrls: ['./new-access-key-dialog.component.scss']
})
export class NewAccessKeyDialogComponent implements OnInit {

  constructor( public dialogRef: MatDialogRef<NewAccessKeyDialogComponent>,
               @Inject(MAT_DIALOG_DATA)
               public data: AccessKeyData) { }

  ngOnInit() {
  }
}

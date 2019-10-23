import {Component, OnInit} from '@angular/core';
import {MatDialogRef} from "@angular/material";
import {FormBuilder, FormGroup, Validators} from "@angular/forms";

// S3 bucket requirements, check https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-s3-bucket-naming-requirements.html
const NAME_PATTERN = /(?!^(\d+\.)+\d+$)(^(([a-z0-9]|[a-z0-9][a-z0-9\-]*[a-z0-9])\.)*([a-z0-9]|[a-z0-9][a-z0-9\-]*[a-z0-9])$)/;

@Component({
  selector: 'app-new-bucket-dialog',
  templateUrl: './new-bucket-dialog.component.html',
  styleUrls: ['./new-bucket-dialog.component.scss']
})
export class NewBucketDialogComponent implements OnInit {
  form: FormGroup;
  private formSubmitAttempt: boolean;

  constructor(private formBuilder: FormBuilder,
              public dialogRef: MatDialogRef<NewBucketDialogComponent>) {
  }

  ngOnInit() {
    this.form = this.formBuilder.group({
      'name': ['', [Validators.required, Validators.minLength(3), Validators.maxLength(63), Validators.pattern(NAME_PATTERN)]],
    });
  }

  get name() {
    return this.form.get('name').value;
  }

  onSubmit(): void {
    if (this.form.valid) {
      this.dialogRef.close(this.name)
    }
    this.formSubmitAttempt = true;
  }

  isFieldInvalid(field: string) {
    return (
      (!this.form.get(field).valid && this.form.get(field).touched) ||
      (this.form.get(field).untouched && this.formSubmitAttempt)
    );
  }

  getErrorMessage(fieldName: string) {
    let field = this.form.get(fieldName)
    if (field.hasError('required')) {
      return 'You must enter a value'
    }
    if (field.hasError('pattern')) {
      return 'Invalid bucket name. Please check bucket naming requirements.'
    }
    if (field.hasError('minlength')) {
      return 'Name must be at least 3 characters long'
    }
    if (field.hasError('maxlength')) {
      return 'Name must be at most 63 characters long'
    }
  }
}

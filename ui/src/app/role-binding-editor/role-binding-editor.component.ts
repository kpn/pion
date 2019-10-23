import {Component, Inject, OnInit} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material";
import {RoleBinding} from "../core/model/role-binding";
import {FormBuilder, FormGroup, Validators} from "@angular/forms";

const VALUE_PATTERN = /^\w[\w\-\_]{3,}$/;

@Component({
  selector: 'app-role-binding-editor',
  templateUrl: './role-binding-editor.component.html',
  styleUrls: ['./role-binding-editor.component.scss']
})
export class RoleBindingEditorComponent implements OnInit {
  roleTypes: string[] = ['readonly', 'user', 'editor', 'admin'];

  form: FormGroup;
  private formSubmitAttempt: boolean;

  constructor(private formBuilder: FormBuilder,
              public dialogRef: MatDialogRef<RoleBindingEditorComponent>,
              @Inject(MAT_DIALOG_DATA) public data: RoleBinding) {
  }

  ngOnInit() {
    let subject;
    if (this.data.subjects && this.data.subjects.length > 0 && this.data.subjects[0]) {
      subject = this.data.subjects[0]
    } else {
      subject = {value: '', type: 'user'}
    }
    this.form = this.formBuilder.group({
      'name': [this.data.name ? this.data.name : '', [Validators.required, Validators.pattern(VALUE_PATTERN)]],
      'roleRef': [this.data.roleRef ? this.data.roleRef : 'readonly', Validators.required],
      'subjectValue': [subject.value ? subject.value : '', [Validators.required, Validators.pattern(VALUE_PATTERN)]],
      'subjectType': [subject.type ? subject.type : 'user', Validators.required],
    });
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
      return `Value must use alphanumeric, '-' or '_' and has at least 4 characters`
    }
  }

  onSubmit(): void {
    if (this.form.valid) {
      let rb = this.getSubmitObject();
      console.log(rb)
      this.dialogRef.close(rb);
    }
    this.formSubmitAttempt = true;
  }

  private getSubmitObject() {
    return {
      name: this.form.controls['name'].value,
      roleRef: this.form.controls['roleRef'].value,
      subjects: [
        {
          value: this.form.controls['subjectValue'].value,
          type: this.form.controls['subjectType'].value,
        }
      ]
    }
  }
}

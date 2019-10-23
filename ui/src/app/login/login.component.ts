import { Component, OnInit } from '@angular/core';
import {FormGroup, FormBuilder, Validators} from "@angular/forms";
import {AuthnService} from "../core/services/authn/authn.service";

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {
  form: FormGroup;
  private formSubmitAttempt: boolean;
  loginFailed: boolean;

  constructor(private fb: FormBuilder,
              private authnService: AuthnService
  ) { }

  ngOnInit() {
    this.form = this.fb.group({     // {5}
      customer: ['', Validators.required],
      username: ['', Validators.required],
      password: ['', Validators.required]
    });
  }

  isFieldInvalid(field: string) { // {6}
    return (
      (!this.form.get(field).valid && this.form.get(field).touched) ||
      (this.form.get(field).untouched && this.formSubmitAttempt)
    );
  }

  onSubmit(): void {
    if (this.form.valid) {
      this.authnService.login(this.form.value)
        .subscribe(resp=>this.onLoginSuccess(resp), resp=> this.onLoginFailure(resp))
    }
    this.formSubmitAttempt = true;
  }

  private onLoginSuccess(resp) {
    console.log(resp);
  }

  private onLoginFailure(resp) {
    console.log(resp);
    this.loginFailed = true
  }

  onFocus() {
    this.loginFailed = false
  }
}

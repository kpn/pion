import {BrowserModule} from '@angular/platform-browser';
import {isDevMode, NgModule} from '@angular/core';
import {DragDropModule} from '@angular/cdk/drag-drop';
import {CdkTableModule} from '@angular/cdk/table';
import {CdkTreeModule} from '@angular/cdk/tree';
import {HttpClientModule} from '@angular/common/http';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';

import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {HttpClientInMemoryWebApiModule} from 'angular-in-memory-web-api';

import {AppComponent} from './app.component';
import {AppRoutingModule} from './app-routing.module';
import {BucketsComponent} from './buckets/buckets.component';
import {AccessKeysComponent} from './accesskeys/access-keys.component';
import {AclEditorDialogComponent} from './acl-editor-dialog/acl-editor-dialog.component';
import {LoginComponent} from './login/login.component';
import {HeaderComponent} from './header/header.component';
import {LoginLayoutComponent} from './login-layout/login-layout.component';
import {HomeLayoutComponent} from './home-layout/home-layout.component';
import {MaterialModule} from "./core/material.module";
import {NewAccessKeyDialogComponent} from './new-access-key-dialog/new-access-key-dialog.component';
import {RoleBindingsComponent} from './role-bindings/role-bindings.component';
import {AuthnService} from "./core/services/authn/authn.service";
import {AuthnGuard} from "./core/services/authn/authn.guard";
import {InMemoryBackendService} from "./core/services/in-memory-backend.service";
import {RolesComponent} from './roles/roles.component';
import {RoleDetailComponent} from './role-detail/role-detail.component';
import {RoleBindingEditorComponent} from './role-binding-editor/role-binding-editor.component';
import {environment} from "../environments/environment";
import { NewBucketDialogComponent } from './new-bucket-dialog/new-bucket-dialog.component';
import {SidenavService} from "./core/services/sidenav.service";

let localModules = [];
if (!environment.remote) {
  localModules = [HttpClientInMemoryWebApiModule.forRoot( //For running local only
    InMemoryBackendService,
    {
      dataEncapsulation: false,
      apiBase: 'restricted/',
    }
  )]
}

@NgModule({
  declarations: [
    AppComponent,
    BucketsComponent,
    AccessKeysComponent,
    AclEditorDialogComponent,
    LoginComponent,
    HeaderComponent,
    LoginLayoutComponent,
    HomeLayoutComponent,
    NewAccessKeyDialogComponent,
    RoleBindingsComponent,
    RolesComponent,
    RoleDetailComponent,
    RoleBindingEditorComponent,
    NewBucketDialogComponent
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    FormsModule,
    HttpClientModule,
    ReactiveFormsModule,
    CdkTableModule,
    CdkTreeModule,
    DragDropModule,
    AppRoutingModule,
    MaterialModule,
    ...localModules
],
entryComponents: [
  AclEditorDialogComponent,
  NewAccessKeyDialogComponent,
  RoleBindingEditorComponent,
  NewBucketDialogComponent
],
  providers
:
[AuthnService, AuthnGuard, SidenavService],
  bootstrap
:
[AppComponent]
})

export class AppModule {
}

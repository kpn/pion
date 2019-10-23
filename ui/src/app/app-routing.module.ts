import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {AccessKeysComponent} from './accesskeys/access-keys.component';
import {BucketsComponent} from './buckets/buckets.component';
import {LoginComponent} from "./login/login.component";
import {LoginLayoutComponent} from "./login-layout/login-layout.component";
import {HomeLayoutComponent} from "./home-layout/home-layout.component";
import {RoleBindingsComponent} from "./role-bindings/role-bindings.component";
import {AuthnGuard} from "./core/services/authn/authn.guard";
import {RolesComponent} from "./roles/roles.component";
import {RoleDetailComponent} from "./role-detail/role-detail.component";

const routes: Routes = [
  {path: '', redirectTo: '/login', pathMatch: 'full'},
  {
    path: 'login', component: LoginLayoutComponent,
    children: [
      {path: '', component: LoginComponent}
    ]
  },
  {
    path: 'main', component: HomeLayoutComponent,
    children: [
      {path: '', redirectTo: 'access-keys', pathMatch: 'full'},
      {path: 'access-keys', component: AccessKeysComponent, canActivate: [AuthnGuard]},
      {path: 'buckets', component: BucketsComponent, canActivate: [AuthnGuard]},
      {path: 'roles', component: RolesComponent, canActivate: [AuthnGuard]},
      {path: 'roles/:name', component: RoleDetailComponent, canActivate: [AuthnGuard]},
      {path: 'role-bindings', component: RoleBindingsComponent, canActivate: [AuthnGuard]},
    ]
  }
];

@NgModule({
  imports: [RouterModule.forRoot(
    routes,
    {useHash: false} // debugging mode
  )
  ],
  exports: [RouterModule]
})
export class AppRoutingModule {
}

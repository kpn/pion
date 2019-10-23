import { Component, OnInit } from '@angular/core';
import {Observable} from "rxjs";
import {AuthnService} from "../core/services/authn/authn.service";
import {SidenavService} from "../core/services/sidenav.service";

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent implements OnInit {

  isLoggedIn$: Observable<boolean>;

  constructor(private authnService: AuthnService, private sidenav: SidenavService) { }

  ngOnInit() {
    this.isLoggedIn$ = this.authnService.isLoggedIn;
  }

  onLogout() {
    this.authnService.logout()
      .subscribe(
        ()=>console.log('logged out'),
        ()=> alert('log-out failed')
      );
  }

  get currentUser() {
    return this.authnService.currentUser;
  }

  toggleSideNav() {
    this.sidenav.toggle()
  }
}

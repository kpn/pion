import {Component, OnInit, ViewChild} from '@angular/core';
import {MatSidenav} from "@angular/material";
import {SidenavService} from "../core/services/sidenav.service";

@Component({
  selector: 'app-home-layout',
  templateUrl: './home-layout.component.html',
  styleUrls: ['./home-layout.component.scss']
})
export class HomeLayoutComponent implements OnInit {

  @ViewChild('snav')
  public sidenav: MatSidenav;

  constructor(private sidenavService: SidenavService) { }

  ngOnInit() {
    this.sidenavService.setSidenav(this.sidenav);
  }

  toggle() {
    return this.sidenav.toggle();
  }

}

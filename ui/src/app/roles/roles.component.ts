import {Component, OnInit, ViewChild} from '@angular/core';
import {MatSort, MatTableDataSource} from "@angular/material";
import {Role} from "../core/model/role";
import {RoleService} from "../core/services/role.service";

@Component({
  selector: 'app-roles',
  templateUrl: './roles.component.html',
  styleUrls: ['./roles.component.scss']
})
export class RolesComponent implements OnInit {
  displayedColumns: string[] = ['state', 'name', 'createdAt'];
  dsRoles: MatTableDataSource<Role>;

  @ViewChild(MatSort) sort: MatSort;

  constructor(private roleService: RoleService) { }

  ngOnInit() {
    this.getRoles();
  }

  private getRoles() {
    this.roleService.getRoles()
      .subscribe(roles => this.updateDataSource(roles),
        err => console.log(`fetching roles failed: ${err}`));
  }

  private updateDataSource(roles: Role[]) {
    this.dsRoles = new MatTableDataSource<Role>(roles ? roles : []);
    this.dsRoles.sort = this.sort;
  }
}

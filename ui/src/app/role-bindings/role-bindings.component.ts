import {Component, OnInit, ViewChild} from '@angular/core';
import {MatDialog, MatSort, MatTableDataSource} from "@angular/material";
import {RoleBinding} from "../core/model/role-binding";
import {RoleBindingService} from "../core/services/role-binding.service";
import {RoleBindingEditorComponent} from "../role-binding-editor/role-binding-editor.component";

@Component({
  selector: 'app-role-bindings',
  templateUrl: './role-bindings.component.html',
  styleUrls: ['./role-bindings.component.scss']
})
export class RoleBindingsComponent implements OnInit {
  displayedColumns: string[] = ['name', 'subjects', 'role', 'createdAt', 'actions'];
  dsRoleBindings: MatTableDataSource<RoleBinding>;

  @ViewChild(MatSort) sort: MatSort;

  constructor(private dialog: MatDialog,
              private roleBindingService: RoleBindingService) {
  }

  ngOnInit() {
    this.getRoleBindings()
  }

  getRoleBindings(): void {
    this.roleBindingService.getRoleBindings()
      .subscribe(bindings => this.updateDataSource(bindings),
        err => this.log('fetching access-keys failed:', err));
  }

  private updateDataSource(bindings) {
    this.dsRoleBindings = new MatTableDataSource(bindings ? bindings : []);
    this.dsRoleBindings.sort = this.sort;
  }

  multipleSubjects(element: RoleBinding): boolean {
    return element.subjects.length > 1;
  }

  newRoleBinding() {
    const dialogRef = this.dialog.open(RoleBindingEditorComponent, {
      data: {}
    });

    dialogRef.afterClosed().subscribe(result => {
      if (!result) {
        return;
      }
      this.log('Saving object', result);
      this.roleBindingService.createRoleBinding(result)
        .subscribe(
          result => {
            console.log(`Saved ${result}`);
            this.dsRoleBindings.data.push(result);
            this.dsRoleBindings.data = this.dsRoleBindings.data.slice();
          },
          err => {
            this.log('failed to save:', err);
          }
        );
    });
  }

  deleteRoleBinding(roleBinding: RoleBinding) {
    if (confirm(`Are you sure to delete the role-binding '${roleBinding.name}'?`)) {
      this.roleBindingService.deleteRoleBinding(roleBinding.name)
        .subscribe(() => this.dsRoleBindings.data = this.dsRoleBindings.data.filter(rb => rb !== roleBinding));
    }

    this.roleBindingService.deleteRoleBinding(roleBinding.name)
  }

  openEditorDialog(roleBinding: RoleBinding) {
    const dialogRef = this.dialog.open(RoleBindingEditorComponent, {
      data: roleBinding,
    });

    dialogRef.afterClosed().subscribe(
      result => {
        if (!result) {
          return;
        }
        this.log('Updating object', result);

        this.roleBindingService.updateRoleBinding(result)
          .subscribe(
            result => this.updateRoleBinding(roleBinding.name, result),
            err => this.log('failed to update:', err)
          );
      }
    );
  }

  log(message: string, object: any) {
    console.log(message);
    console.log(object);
  }

  private updateRoleBinding(name: string, updatedRB: RoleBinding) {
    let target = this.dsRoleBindings.data.find(b => b.name === name);
    // update object attributes without loosing reference
    Object.assign(target, updatedRB);
  }
}

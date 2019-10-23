import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { RoleBindingEditorComponent } from './role-binding-editor.component';

describe('RoleBindingEditorComponent', () => {
  let component: RoleBindingEditorComponent;
  let fixture: ComponentFixture<RoleBindingEditorComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ RoleBindingEditorComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(RoleBindingEditorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

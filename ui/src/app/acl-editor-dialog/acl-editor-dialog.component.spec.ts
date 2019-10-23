import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AclEditorDialogComponent } from './acl-editor-dialog.component';

describe('AclEditorDialogComponent', () => {
  let component: AclEditorDialogComponent;
  let fixture: ComponentFixture<AclEditorDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AclEditorDialogComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AclEditorDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

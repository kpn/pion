import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NewAccessKeyDialogComponent } from './new-access-key-dialog.component';

describe('NewAccessKeyDialogComponent', () => {
  let component: NewAccessKeyDialogComponent;
  let fixture: ComponentFixture<NewAccessKeyDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NewAccessKeyDialogComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NewAccessKeyDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

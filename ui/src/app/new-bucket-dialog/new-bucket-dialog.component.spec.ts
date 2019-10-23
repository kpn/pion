import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NewBucketDialogComponent } from './new-bucket-dialog.component';

describe('NewBucketDialogComponent', () => {
  let component: NewBucketDialogComponent;
  let fixture: ComponentFixture<NewBucketDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NewBucketDialogComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NewBucketDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

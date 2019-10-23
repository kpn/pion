import { TestBed } from '@angular/core/testing';

import { BucketService } from './bucket.service';

describe('BucketService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: BucketService = TestBed.get(BucketService);
    expect(service).toBeTruthy();
  });
});

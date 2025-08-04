import { Module } from '@nestjs/common';
import { PharaohService } from './pharaoh.service';
import { PharaohResolver } from './pharaoh.resolver';

@Module({
  providers: [PharaohService, PharaohResolver],
})
export class PharaohModule {}
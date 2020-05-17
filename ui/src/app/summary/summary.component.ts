import { Component, OnInit, Input } from '@angular/core';
import { StatusFilter } from '../todo';

@Component({
    selector: 'app-summary',
    templateUrl: './summary.component.html',
    styleUrls: ['./summary.component.scss'],
})
export class SummaryComponent implements OnInit {
    @Input() todoCount: number;
    @Input() statusFilter: StatusFilter;

    shouldDisplaySummary: boolean;

    ngOnInit(): void {
        this.shouldDisplaySummary =
            this.todoCount > 0 && ['', 'todo'].includes(this.statusFilter);
    }
}

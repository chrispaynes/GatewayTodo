import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { TodosComponent } from './todos/todos.component';


const routes: Routes = [
  { path: 'todo', component: TodosComponent },
  { path: 'completed', component: TodosComponent },
  { path: 'archived', component: TodosComponent },
  { path: 'all', component: TodosComponent },
  { path: '', component: TodosComponent },
  { path: '**', component: TodosComponent }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }

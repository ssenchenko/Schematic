use async_graphql::{Schema, EmptyMutation, EmptySubscription, Error, Object, Result, Context};
use graphql_schema::{Topology, Resource};

struct QueryRoot;

#[Object]
impl QueryRoot {
    async fn topology(
        &self,
        ctx: &Context<'_>,
        type_name: String,
        identifier: String,
    ) -> Result<Topology> {
        Err(Error::new("placeholder"))
    }

    async fn resources(
        &self,
        ctx: &Context<'_>,
        type_name: String,
        next_token: Option<String>,
        resource_model: Option<String>,
    ) -> Result<Vec<Resource>> {
        Err(Error::new("placeholder"))
    }

    async fn resource(
        &self,
        ctx: &Context<'_>,
        type_name: String,
        identifier: String,
    ) -> Result<Resource> {
        Err(Error::new("placeholder"))
    }
}


fn main() {
    let schema = Schema::build(QueryRoot, EmptyMutation, EmptySubscription).finish();
    let sdl = schema.sdl();
    println!("{}", sdl);
}

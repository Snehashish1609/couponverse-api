# Couponverse API

RESTful API to manage and apply different types of discount coupons (cart-wise, product-wise, and BxGy) for an e-commerce platform.

### Run local

`go run cmd/main.go`


## Implementation

CouponVerse is a Go web-server implemented using `gorilla-mux`, which serves as a api service to manage coupons for an e-commerce app. For this project I have used a Postgres Database that enables the ability to add, update, get as well as delete coupons using simple API calls.

Currently I have implemented 3 types of coupons, namely: "cart-wise", "product-wise", and "bxgy".
As the names might suggest:
1. `cart-wise`: Apply a discount to the entire cart if the total amount exceeds a certain threshold.
2. `product-wise`: Apply a discount to specific products.
3. `bxgy`: “Buy X, Get Y” deals with a repetition limit and can be applicable to a set of products (e.g., Buy 3 of Product X or Product Y, get 1 of Product P and Product Q free and so on).

**_Note:_**
*Why use Postgres? Mostly familiarity. But also, as it's substantially fast.*
*Why `gorilla-mux`? Familiarity and ease of use. More customization.*

### API Endpoints:

* POST /coupons: Create a new coupon.
* GET /coupons: Retrieve all coupons.
* GET /coupons/{id}: Retrieve a specific coupon by its ID.
* PUT /coupons/{id}: Update a specific coupon by its ID.
* DELETE /coupons/{id}: Delete a specific coupon by its ID.
* POST /applicable-coupons: Fetch all applicable coupons for a given cart and calculate the total discount that will be applied by each coupon.
* POST /apply-coupon/{id}: Apply a specific coupon to the cart and return the updated cart with discounted prices for each item.

**Code:**

Currently the main code resides in `cmd/main.go`, where we create global config, database client and coupon handler client and pass them on further. All database related operations have been abstracted using a Database Client interface, using GORM under the hood. This ensuures that we can swap off databases anytime, easily.
On the other hand, most of the business logic for handling coupons is in `handlers/v1`, where we use DBClient to talk to the database and apply logic on top of the data.
The code is fairly extensible even though improvements could be done (will list down in `Improvements`). Database client and http handlers have interfaces attached to them, that provides the ability to mock these elements and write independent unit test cases.
To do: dynamic discount calculation per coupon type (more in point 3, Improvements). We can utilize Go interfaces and structs to make the coupon type extensibility more easier and smoother.

**Cases:**
1. User can add as well use coupons of `cart-wise` coupon type. This would apply a discount to the entire cart if total amount exceeds a specified threshold.
2. In `cart-wise` coupon type, a user can specify a discount `Cap`, that would cap the discount to the provided amount (if supplied)
3. User can add and use coupons of `product-wise` coupon type. This would apply to specific products in cart.
4. In `product-wise` coupon type, a user can also specify a discount `Cap`, that would cap the discount to the provided amount (if supplied)
5. BuyXGetY or `bxgy` type coupons can be used by the application owners to create/define coupons that would offer discounts based on combination of products in cart.
6. `bxgy` type coupons can be used to add a repetition limit, such that we can limit how many times a product can be discounted/combined as free.
7. After checkout, the `/applicable-coupons` and `apply-coupons` endpoints can be incorporated in any e-commerce app, which would want to extend to this project (CouponVerse).
8. ...

**Additional Use Cases (not implemented, yet):**
1. Time based coupons: Admins provide when a coupon expires, and business logic could handle if a fetched coupon has expired(delete coupon in that case) or not (apply coupon). This would require addition of expiration timestamp in coupon schema and can be implemented.
2. Coupons based on users: Like a loyalty program, new coupon type can be implemented based on user type. [Not implemented as this would require access to user info along with cart details in payload]
3. Combine coupons: Users could be given the ability to combine or stack coupons on their cart [Not implemented, as it goes beyond the assumption that only 1 coupon can be applied to a single cart]
4. First-come-first-serve coupons: Special coupons could be added with limited quotas. [Not implemented, as this could add another change in the DB schema]
5. ...


### Assumptions

* We want to restrict the user to only be able to apply one coupon for their cart at a time.
* We cannot access User details for the cart payload provided. That is no User info scope.
* This web server is meant to be used by a e-commerce service and not a real user. It only incorporates backend api access and no UI.


### Challenges

* Currently I am storing the coupon `details` as a jsonb object in DB complemented by Go string. So when calling the APIs, the details portion is now to be passed as a string. Need to figure out a better solution. Using a NoSQL database would be easier, as details schema changes for each coupon type.


### Improvements

* Using MongoDB instead of Postgres might be beneficial if we want more flexibility in coupon structures. It should be easy to do
* Better error handling.
* implementing an interface (maybe DiscountCalulator) to abstract out the discount calculation for different coupon types. Currently it has been handled using different functions and a simple `switch`.
* optimize calculation logic for `BxGy` coupon type. Currently I've used brute force.
* adding authentication barrier for coupon CRUD operations could be done (using an auth middleware). This would ensure users can only calculate and apply discounts and not anything else.
* Using gRPC over REST. We could use gRPC as this is not a user facing service and only is used by another microservice. gRPC might give more speed and consistency, but would also require fairly static schemas.
* More refactoring :)


### ToDos

- use interface for coupon type?
- unit tests
- better database schema?
- better way to handle details per coupon type?
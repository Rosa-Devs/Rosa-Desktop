export namespace manifest {
	
	export class Manifest {
	    name: string;
	    uid: string;
	    pubsub: string;
	    chiper: string;
	    optional: string;
	
	    static createFrom(source: any = {}) {
	        return new Manifest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.uid = source["uid"];
	        this.pubsub = source["pubsub"];
	        this.chiper = source["chiper"];
	        this.optional = source["optional"];
	    }
	}

}

export namespace models {
	
	export class ProfileStorePublic {
	    type: string;
	    name: string;
	    id: string;
	    pub: string;
	    avatar: string;
	
	    static createFrom(source: any = {}) {
	        return new ProfileStorePublic(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.name = source["name"];
	        this.id = source["id"];
	        this.pub = source["pub"];
	        this.avatar = source["avatar"];
	    }
	}
	export class Message {
	    datatype: number;
	    sender: ProfileStorePublic;
	    senderid: string;
	    data: string;
	    time: string;
	    sign: string;
	    valid: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.datatype = source["datatype"];
	        this.sender = this.convertValues(source["sender"], ProfileStorePublic);
	        this.senderid = source["senderid"];
	        this.data = source["data"];
	        this.time = source["time"];
	        this.sign = source["sign"];
	        this.valid = source["valid"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}


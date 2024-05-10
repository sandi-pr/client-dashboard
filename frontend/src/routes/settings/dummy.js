export const datas = [
  {
    id: 1,
    name: "User 01",
    role: "admin",
    email: "user01@gmail.com",
    actions: {
      objects: [
        {
          name: "Delete",
          href: "/companies/delete/1454523",
          iconOnly: true,
        },
        {
          name: "Edit",
          href: "/companies/update/1454523",
          iconOnly: true,
        },
      ],
    },
  },
  {
    id: 2,
    name: "User 02",
    role: "admin",
    email: "user02@gmail.com",
    actions: {
      objects: [
        {
          name: "Delete",
          href: "/companies/delete/1454523",
          iconOnly: true,
        },
        {
          name: "Edit",
          href: "/companies/update/1454523",
          iconOnly: true,
        },
      ],
    },
  },
];

export const informationFormDatas = {
  grid: 4,
  title: "My Details",
  titleGrid: 1,
  innerGrid: 5,
  data: [
    {
      name: "Full Name",
      value: "Admin",
      grid: 2,
    },
    {
      name: "Username",
      value: "User 01",
      grid: 2,
    },
    {
      name: "User ID",
      value: 1,
      grid: 1,
    },
    {
      name: "Email",
      value: "user01@gmail.com",
      grid: 5,
    },
    // {
    //   name: "Role",
    //   value: "Admin",
    //   grid: 2,
    // },
  ],
};

export const configs = {
  headers: "Admin List",
  actions: [
    {
      name: "Add",
      icon: "plus",
      href: "/",
      type: "brand",
    },
  ],
};

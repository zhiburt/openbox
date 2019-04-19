export default function() {
  return [
    {
      title: "Reports",
      htmlBefore: '<i class="material-icons">table_chart</i>',
      to: "/tables",
    },
    {
      title: "Files",
      htmlBefore: '<i class="material-icons">folder</i>',
      to: "/blog-posts",
    },
    {
      title: "User Profile",
      htmlBefore: '<i class="material-icons">person</i>',
      to: "/user-profile-lite",
    },
    {
      title: "Coins",
      htmlBefore: '<i class="material-icons">error</i>',
      to: "/errors",
    }
  ];
}
